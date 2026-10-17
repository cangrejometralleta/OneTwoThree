package com.example.campus;

import com.example.app.EnrollmentNotifier;
import com.example.app.RutRegistry;
import com.example.codec.Json;
import com.example.school.Rut;
import com.example.school.Student;

import java.io.IOException;
import java.lang.System.Logger;
import java.lang.System.Logger.Level;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

/**
 * Office Answers the Registry and Announces Enrolments: one Adapter, two Ports.
 * The Registry is a JSON File of the Ruts that are Real; the Notice goes to the Log.
 */
public final class Office implements RutRegistry, EnrollmentNotifier {

  private final Set<Rut> registered;
  private final Logger logger;

  private Office(Set<Rut> registered, Logger logger) {
    this.registered = Set.copyOf(registered);
    this.logger = logger;
  }

  /**
   * openOffice Reads the Registry File once, at Startup.
   * A File that is Missing, Empty, not an Array, Holds a second Value or a Text that
   * is not a valid Rut Stops Startup.
   */
  public static Office openOffice(Path path, Logger logger) throws IOException {
    Object read = Json.readJsonText(Files.readAllBytes(path))
        .orElseThrow(() -> new IOException("registry file " + path + " is not one JSON value"));
    if (!(read instanceof List<?> items)) {
      throw new IOException("registry file " + path + " Must Hold an Array");
    }

    Set<Rut> ruts = new HashSet<>();
    for (Object item : items) {
      if (!(item instanceof String text) || !new Rut(text).looksValid()) {
        throw new IOException("registry file " + path + ": " + item + " is not a valid RUT");
      }
      ruts.add(new Rut(text));
    }

    return new Office(ruts, logger);
  }

  /** confirmRut Answers whether the Registry Knows a Rut. This one never Fails to Answer. */
  @Override
  public boolean confirmRut(Rut rut) {
    return registered.contains(rut);
  }

  /** announceEnrollment Logs the Notice, Naming the Student by Id and never by Rut or Name. */
  @Override
  public void announceEnrollment(Student student) {
    logger.log(Level.INFO, "student " + student.id().value() + " is enrolled");
  }
}
