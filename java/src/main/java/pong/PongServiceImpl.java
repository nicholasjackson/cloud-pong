package pong;

import com.hashicorp.pong.PongData;
import com.hashicorp.pong.PongServiceGrpc;
import io.grpc.stub.StreamObserver;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;


public class PongServiceImpl extends PongServiceGrpc.PongServiceImplBase {
    private static Logger log = LogManager.getLogger(PongServiceImpl.class.getName());
    private StreamObserver<PongData> streamObserver;

    @Override
    public StreamObserver<PongData> serverStream (StreamObserver<PongData> responseObserver) {
        return createStreamObserver(responseObserver);
    }

    StreamObserver<PongData> createStreamObserver(StreamObserver<PongData> responseObserver) {
        return this.streamObserver = new StreamObserver<>() {
            @Override
            public void onNext(PongData value) {
                responseObserver.onNext(value);
                log.info("got server event BatX: {}, BatY: {}, BallX: {}, BallY: {}",
                        value.getBat().getX(),
                        value.getBat().getY(),
                        value.getBall().getX(),
                        value.getBall().getY());
            }

            @Override
            public void onError(Throwable t) {
                log.error("error reading server stream : {}", t.getMessage());
            }

            @Override
            public void onCompleted() {
                log.debug("completed server stream");
            }
        };
    }

    @Override
    public StreamObserver<PongData> clientStream (StreamObserver<PongData> responseObserver) {
        return new StreamObserver<>() {
            @Override
            public void onNext(PongData value) {
                responseObserver.onNext(value);
                streamObserver.onNext(value);
                log.info("got client event BatX: {}, BatY: {}, BallX: {}, BallY: {}",
                        value.getBat().getX(),
                        value.getBat().getY(),
                        value.getBall().getX(),
                        value.getBall().getY());
            }

            @Override
            public void onError(Throwable t) {
                log.error("error client stream : {}", t.getMessage());
            }

            @Override
            public void onCompleted() {
                log.debug("completed client stream");
            }
        };
    }
}
