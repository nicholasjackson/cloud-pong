package pong;

import io.grpc.*;
import org.apache.logging.log4j.Level;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.logging.log4j.core.config.Configurator;

public class ApiServer {
    private final static Logger log = LogManager.getLogger(ApiServer.class.getName());
    private static int player = Integer.parseInt(getEnvOrDefault("PLAYER", "1"));
    private static int port = Integer.parseInt(getEnvOrDefault("PORT", "6000"));
    private static String upstream = getEnvOrDefault("UPSTREAM_ADDRESS", "localhost:6001");

    public static void main( String[] args ) throws Exception {

        // create a gRPC client to communicate with the other server
        //

        Configurator.setRootLevel(Level.INFO);
        Server server = ServerBuilder
                .forPort(port)
                .addService(new PongServiceImpl()).build();
        log.info("Listening on port {}, player {}", port, player);
        server.start();

        server.awaitTermination();

    }

    public static String getEnvOrDefault(String env, String defaultValue) {
      String value = System.getenv(env);

        if (value == null) {
          return defaultValue;
        }

      return value;
    }
}
