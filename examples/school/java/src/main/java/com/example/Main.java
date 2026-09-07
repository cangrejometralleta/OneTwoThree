package com.example;

import com.example.api.SchoolAPI;
import com.example.serving.HTTPServer;
import com.example.settings.SchoolConfig;
import com.example.settings.Settings;
import com.example.store.CourseRows;
import com.example.store.School;
import com.example.store.StudentRows;
import com.example.tokens.AccessTokens;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.web.servlet.function.RouterFunction;
import org.springframework.web.servlet.function.ServerResponse;

/**
 * main Casts the Players, then Steps off the Stage.
 * It is the only Class that Knows every Package by Name.
 */
@SpringBootApplication
public class Main {

  public static void main(String[] args) {
    SchoolConfig config = SchoolConfig.loadSchoolConfig(Settings.findSchoolDataRoot(),
        Settings.readProcessEnvironment());

    System.setProperty("server.port", String.valueOf(config.port()));
    System.setProperty("school.token.secret", config.tokenSecret());
    System.setProperty("school.token.life", String.valueOf(config.tokenLifeSeconds()));
    System.setProperty("school.database.path", config.databasePath());

    SpringApplication.run(Main.class, args);
  }

  @Bean
  RouterFunction<ServerResponse> schoolRoutes(StudentRows studentRows, CourseRows courseRows) {
    School registry = new School(studentRows, courseRows);
    SchoolAPI service = new SchoolAPI(registry, registry,
        AccessTokens.buildAccessTokens(System.getProperty("school.token.secret", ""),
            Integer.parseInt(System.getProperty("school.token.life", "600"))));

    return new HTTPServer(service.declareSchoolRoutes()).serveRoutes();
  }
}
