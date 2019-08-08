package pong;

import io.grpc.*;

import java.io.IOException;
import java.util.logging.Level;
import java.util.logging.Logger;


//func main() {
//        env.Parse()
//
//        logger = hclog.Default()
//        apiClient := client.New(*upstream)
//        apiClient.DialAsync()
//
//        grpcServer := grpc.NewServer()
//        server := server.New(logger, apiClient)
//        pb.RegisterPongServiceServer(grpcServer, server)
//
//        l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
//        if err != nil {
//        logger.Error("Failed to listen", "error", err)
//        os.Exit(1)
//        }
//
//        logger.Info("Listening on port", "port", *port, "player", *player)
//        grpcServer.Serve(l)
//        }
//}

public class ApiServer {
    private final static Logger LOGGER = Logger.getLogger(ApiServer.class.getName());

    private static void setup() throws IOException {
        LOGGER.setLevel(Level.INFO);
    }

    public static void main( String[] args ) throws Exception {
        setup();
        Server server = ServerBuilder
                .forPort(8080)
                .addService(new PongServiceImpl()).build();

        server.start();
        server.awaitTermination();
    }
}
