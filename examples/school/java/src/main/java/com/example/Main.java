package com.example;

import com.example.app.SchoolService;
import com.example.campus.Office;
import com.example.handler.SchoolHandler;
import com.example.school.SchoolRules;
import com.example.serving.spring.SpringServer;
import com.example.serving.stdlib.StdlibServer;
import com.example.settings.GlobalConstants;
import com.example.settings.SchoolConfig;
import com.example.settings.Settings;
import com.example.store.SchoolStore;
import com.example.tokens.AccessTokens;
import com.example.transport.Server;

import java.nio.file.Path;
import java.time.Clock;
import java.util.Map;
import java.util.function.Supplier;

/**
 * Main Casts the Players, then Steps off the Stage.
 * It is the only Class that Knows every Package by Name.
 */
public final class Main {

  /** SERVER_ADAPTERS is every Framework this Runtime Implements; the Loader only Receives the Names. */
  static final Map<String, Supplier<Server>> SERVER_ADAPTERS = Map.of(
      "stdlib", StdlibServer::new,
      "spring", SpringServer::new);

  private Main() {
  }

  public static void main(String[] args) {
    try {
      startSchool();
    } catch (Exception refused) {
      System.err.println("❌ School Refused to Start: " + refused.getMessage());
      System.exit(1);
    }
  }

  private static void startSchool() throws Exception {
    Path root = Settings.findSchoolDataRoot(Path.of(""));
    GlobalConstants constants = Settings.readGlobalConstants(root);
    SchoolConfig config = SchoolConfig.loadSchoolConfig(root, System.getenv(), SERVER_ADAPTERS.keySet());

    SchoolStore store = SchoolStore.openSchoolStore(config.databasePath());
    Office office = Office.openOffice(config.registryPath(), System.getLogger("school.campus"));
    SchoolService service = new SchoolService(new SchoolRules(constants.minimumAgeYears()), store, store,
        office, office, System.getLogger("school.app"));
    SchoolHandler door = new SchoolHandler(service,
        AccessTokens.buildAccessTokens(config.tokenSecret(), config.tokenLifeSeconds(), Clock.systemUTC()));

    Server server = SERVER_ADAPTERS.get(config.serverAdapter()).get();
    Server.Running running = server.serveRoutes(door.declareSchoolRoutes(), config.port());
    System.out.println("✅ School Listening on :" + running.port() + " through " + config.serverAdapter());
  }
}
