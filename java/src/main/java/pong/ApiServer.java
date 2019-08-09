package pong;

import io.grpc.*;
import org.apache.logging.log4j.Level;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.logging.log4j.core.config.Configurator;

public class ApiServer {
    private final static Logger log = LogManager.getLogger(ApiServer.class.getName());
    private static int player = Integer.parseInt(getEnvOrDefault("PLAYER", "2"));
    private static int port = Integer.parseInt(getEnvOrDefault("PORT", "6001"));
    private static String upstream = getEnvOrDefault("UPSTREAM_ADDRESS", "localhost:6000");

    public static void main( String[] args ) throws Exception {
        Configurator.setRootLevel(Level.INFO);
        PongClient client = new PongClient(upstream);
        PongServiceImpl pongService = new PongServiceImpl();
        pongService.setOtherGameServerClient(client);
        Server server = ServerBuilder
                .forPort(port)
                .addService(pongService).build();


        try {
            server.start();
            log.info("Listening on port {}, player {}", port, player);
        } finally {
            server.awaitTermination();
            client.shutdown();
        }

    }

    public static String getEnvOrDefault(String env, String defaultValue) {
      String value = System.getenv(env);

        if (value == null) {
          return defaultValue;
        }

      return value;
    }
}
