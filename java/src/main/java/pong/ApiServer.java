package pong;

import io.grpc.*;
import org.apache.logging.log4j.Level;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.logging.log4j.core.config.Configurator;

public class ApiServer {
    private final static Logger log = LogManager.getLogger(ApiServer.class.getName());

    private static void setup(){
        Configurator.setRootLevel(Level.INFO);
    }

    public static void main( String[] args ) throws Exception {
        setup();
        Server server = ServerBuilder
                .forPort(6000)
                .addService(new PongServiceImpl()).build();
        log.info("started server on port {}", 6000);

        server.start();
        server.awaitTermination();
    }
}
