package com.example.campus;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

import com.example.app.EnrollmentNotifier;
import com.example.app.RutRegistry;
import com.example.school.CourseId;
import com.example.school.FullName;
import com.example.school.Rut;
import com.example.school.Student;
import com.example.school.StudentId;

import java.io.IOException;
import java.lang.System.Logger;
import java.nio.file.Files;
import java.nio.file.Path;
import java.text.MessageFormat;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.ResourceBundle;
import java.util.function.Supplier;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

class OfficeTest {

  @TempDir
  Path folder;

  /** Notebook is a Logger that Keeps every Line it is Told. */
  static final class Notebook implements Logger {
    final List<String> lines = new ArrayList<>();

    @Override
    public String getName() {
      return "notebook";
    }

    @Override
    public boolean isLoggable(Level level) {
      return true;
    }

    @Override
    public void log(Level level, ResourceBundle bundle, String message, Throwable thrown) {
      lines.add(message);
    }

    @Override
    public void log(Level level, ResourceBundle bundle, String format, Object... params) {
      lines.add(MessageFormat.format(format, params));
    }

    @Override
    public void log(Level level, Supplier<String> message) {
      lines.add(message.get());
    }
  }

  private Path writeRegistry(String name, String content) throws IOException {
    return Files.writeString(folder.resolve(name), content);
  }

  @Test
  void oneOfficeFillsBothPorts() throws IOException {
    Office office = Office.openOffice(writeRegistry("registry.json", "[]"), new Notebook());

    assertTrue(office instanceof RutRegistry && office instanceof EnrollmentNotifier);
  }

  @Test
  void confirmRutKnowsOnlyWhatTheFileHoldsInLowerCase() throws IOException {
    Office office = Office.openOffice(writeRegistry("registry.json", "[\"11111111-1\",\"10000013-K\",\"11111111-1\"]"),
        new Notebook());

    assertTrue(office.confirmRut(new Rut("11111111-1")));
    assertTrue(office.confirmRut(new Rut("10000013-k")));
    assertFalse(office.confirmRut(new Rut("22222222-2")));
  }

  @Test
  void openOfficeStopsOnWhatItCannotTrust() throws IOException {
    Map<String, String> broken = Map.of(
        "empty", "",
        "broken json", "[\"11111111-1\"",
        "not an array", "{\"ruts\": []}",
        "a second value", "[] []",
        "a bad check", "[\"11111111-9\"]",
        "a non rut text", "[\"hello\"]",
        "not a text", "[11111111]");

    for (Map.Entry<String, String> file : broken.entrySet()) {
      Path path = writeRegistry(file.getKey().replace(' ', '-') + ".json", file.getValue());
      assertThrows(IOException.class, () -> Office.openOffice(path, new Notebook()), file.getKey());
    }
    assertThrows(IOException.class, () -> Office.openOffice(folder.resolve("absent.json"), new Notebook()));
  }

  @Test
  void theShippedRegistryHoldsOnlyRealRuts() throws IOException {
    Office office = Office.openOffice(Path.of("..", "config", "registry.json"), new Notebook());

    assertTrue(office.confirmRut(new Rut("12345678-5")));
  }

  @Test
  void anAnnouncementNamesTheStudentByIdOnly() throws IOException {
    Notebook notebook = new Notebook();
    Office office = Office.openOffice(writeRegistry("registry.json", "[]"), notebook);

    office.announceEnrollment(new Student(new StudentId(7), new Rut("11111111-1"), new FullName("Ana"), 20,
        new CourseId(1)));

    assertTrue(notebook.lines.equals(List.of("student 7 is enrolled")), notebook.lines.toString());
  }
}
