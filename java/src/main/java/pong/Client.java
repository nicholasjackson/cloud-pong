package pong;

import com.hashicorp.pong.PongServiceGrpc;
import io.grpc.*;
import io.grpc.stub.*;
import com.hashicorp.pong.PongServiceGrpc.*;
import com.hashicorp.pong.PongData;
import com.hashicorp.pong.Ball;
import com.hashicorp.pong.Bat;

public class Client {
    public static void main( String[] args ) throws Exception {
        // Channel is the abstraction to connect to a service endpoint
        // Let's use plaintext communication because we don't have certs
        final ManagedChannel channel = ManagedChannelBuilder.forTarget("localhost:8080")
                .usePlaintext()
                .build();

        // It is up to the client to determine whether to block the call
        // Here we create a blocking stub, but an async stub,
        // or an async stub with Future are always possible.
        PongServiceStub stub = PongServiceGrpc.newStub(channel);

        Ball ball = Ball.newBuilder().setX(100).setY(100).build();
        Bat bat = Bat.newBuilder().setX(10).setY(50).build();
        PongData data = PongData.newBuilder()
                .setBall(ball)
                .setBat(bat)
                .setHit(true)
                .build();

        StreamObserver<PongData> toServer = new StreamObserver<>() {
            public void onNext(PongData response) {
                System.out.println(response);
            }
            public void onError(Throwable t) {
            }
            public void onCompleted() {
                channel.shutdownNow();
            }
        };

        toServer.onNext(data);
        toServer.onCompleted();
    }
}
