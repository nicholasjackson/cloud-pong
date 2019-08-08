package pong;

import io.grpc.*;
import org.apache.logging.log4j.Level;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.logging.log4j.core.config.Configurator;

public class ApiServer {
    private final static Logger log = LogManager.getLogger(ApiServer.class.getName());
    private static int player = Integer.parseInt(System.getProperty("PLAYER", "1"));
    private static int port = Integer.parseInt(System.getProperty("PORT", "6000"));
    private static String upstream = System.getProperty("UPSTREAM_ADDRESS", "localhost:6001");

    public static void main( String[] args ) throws Exception {
        Configurator.setRootLevel(Level.INFO);
        Server server = ServerBuilder
                .forPort(port)
                .addService(new PongServiceImpl()).build();
        log.info("Listening on port {}, player {}", port, player);
        server.start();

        server.awaitTermination();

    }
}
