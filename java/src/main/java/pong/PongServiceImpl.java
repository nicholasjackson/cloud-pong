package pong;

import com.hashicorp.pong.PongData;
import com.hashicorp.pong.PongServiceGrpc;
import io.grpc.stub.StreamObserver;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;


public class PongServiceImpl extends PongServiceGrpc.PongServiceImplBase {
    private static Logger log = LogManager.getLogger(PongServiceImpl.class.getName());

    @Override
    public StreamObserver<PongData> serverStream (StreamObserver<PongData> responseObserver) {
        return new StreamObserver<>() {
            @Override
            public void onNext(PongData value) {
                responseObserver.onNext(PongData.getDefaultInstance());
                log.info("got event BatX: {}, BatY: {}, BallX: {}, BallY: {}",
                        PongData.getDefaultInstance().getBat().getX(),
                        PongData.getDefaultInstance().getBat().getY(),
                        PongData.getDefaultInstance().getBall().getX(),
                        PongData.getDefaultInstance().getBall().getY());
            }

            @Override
            public void onError(Throwable t) {
                log.error("error reading stream : {}", t.getMessage());
            }

            @Override
            public void onCompleted() {
                log.debug("completed stream");
            }
        };
    }

    @Override
    public StreamObserver<PongData> clientStream (StreamObserver<PongData> responseObserver) {
        return new StreamObserver<>() {
            @Override
            public void onNext(PongData value) {
                responseObserver.onNext(PongData.getDefaultInstance());
                log.info("got event BatX: {}, BatY: {}, BallX: {}, BallY: {}",
                        PongData.getDefaultInstance().getBat().getX(),
                        PongData.getDefaultInstance().getBat().getY(),
                        PongData.getDefaultInstance().getBall().getX(),
                        PongData.getDefaultInstance().getBall().getY());
            }

            @Override
            public void onError(Throwable t) {
                log.error("error reading stream : {}", t.getMessage());
            }

            @Override
            public void onCompleted() {
                log.debug("completed stream");
            }
        };
    }
}
