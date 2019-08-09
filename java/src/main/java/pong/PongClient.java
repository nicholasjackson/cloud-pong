package pong;

import com.hashicorp.pong.PongData;
import com.hashicorp.pong.PongServiceGrpc;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.stub.StreamObserver;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.util.concurrent.TimeUnit;

public class PongClient {
    private static Logger log = LogManager.getLogger(PongClient.class.getName());

    private final ManagedChannel channel;
    private final PongServiceGrpc.PongServiceStub asyncStub;

    public PongClient(String target) {
        this(ManagedChannelBuilder.forTarget(target).usePlaintext());
    }

    public PongClient(ManagedChannelBuilder<?> channelBuilder) {
        channel = channelBuilder.build();
        this.asyncStub = PongServiceGrpc.newStub(channel);
    }

    public void shutdown() throws InterruptedException {
        channel.shutdown().awaitTermination(5, TimeUnit.SECONDS);
    }

    public void send(PongData value) {
        log.info("sending to server stream");
        StreamObserver<PongData> serverStreamObserver = this.asyncStub.serverStream(new StreamObserver<>() {
            @Override
            public void onNext(PongData value) {
                log.info("server stream BatX: {}, BatY: {}, BallX: {}, BallY: {}",
                        value.getBat().getX(),
                        value.getBat().getY(),
                        value.getBall().getX(),
                        value.getBall().getY());
            }

            @Override
            public void onError(Throwable t) {
                log.error("error server stream : {}", t.getMessage());
            }

            @Override
            public void onCompleted() {
                log.debug("completed server stream");
            }
        });

        try {
            serverStreamObserver.onNext(value);
        } catch (Exception e) {
            serverStreamObserver.onError(e);
            throw e;
        }
        serverStreamObserver.onCompleted();
    }
}
