package com.hashicorp.pong;

import static io.grpc.MethodDescriptor.generateFullMethodName;
import static io.grpc.stub.ClientCalls.asyncBidiStreamingCall;
import static io.grpc.stub.ServerCalls.asyncBidiStreamingCall;
import static io.grpc.stub.ServerCalls.asyncUnimplementedStreamingCall;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.22.1)",
    comments = "Source: api.proto")
public final class PongServiceGrpc {

  private PongServiceGrpc() {}

  public static final String SERVICE_NAME = "pong.PongService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<com.hashicorp.pong.PongData,
      com.hashicorp.pong.PongData> getClientStreamMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "ClientStream",
      requestType = com.hashicorp.pong.PongData.class,
      responseType = com.hashicorp.pong.PongData.class,
      methodType = io.grpc.MethodDescriptor.MethodType.BIDI_STREAMING)
  public static io.grpc.MethodDescriptor<com.hashicorp.pong.PongData,
      com.hashicorp.pong.PongData> getClientStreamMethod() {
    io.grpc.MethodDescriptor<com.hashicorp.pong.PongData, com.hashicorp.pong.PongData> getClientStreamMethod;
    if ((getClientStreamMethod = PongServiceGrpc.getClientStreamMethod) == null) {
      synchronized (PongServiceGrpc.class) {
        if ((getClientStreamMethod = PongServiceGrpc.getClientStreamMethod) == null) {
          PongServiceGrpc.getClientStreamMethod = getClientStreamMethod = 
              io.grpc.MethodDescriptor.<com.hashicorp.pong.PongData, com.hashicorp.pong.PongData>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.BIDI_STREAMING)
              .setFullMethodName(generateFullMethodName(
                  "pong.PongService", "ClientStream"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.hashicorp.pong.PongData.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.hashicorp.pong.PongData.getDefaultInstance()))
                  .setSchemaDescriptor(new PongServiceMethodDescriptorSupplier("ClientStream"))
                  .build();
          }
        }
     }
     return getClientStreamMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.hashicorp.pong.PongData,
      com.hashicorp.pong.PongData> getServerStreamMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "ServerStream",
      requestType = com.hashicorp.pong.PongData.class,
      responseType = com.hashicorp.pong.PongData.class,
      methodType = io.grpc.MethodDescriptor.MethodType.BIDI_STREAMING)
  public static io.grpc.MethodDescriptor<com.hashicorp.pong.PongData,
      com.hashicorp.pong.PongData> getServerStreamMethod() {
    io.grpc.MethodDescriptor<com.hashicorp.pong.PongData, com.hashicorp.pong.PongData> getServerStreamMethod;
    if ((getServerStreamMethod = PongServiceGrpc.getServerStreamMethod) == null) {
      synchronized (PongServiceGrpc.class) {
        if ((getServerStreamMethod = PongServiceGrpc.getServerStreamMethod) == null) {
          PongServiceGrpc.getServerStreamMethod = getServerStreamMethod = 
              io.grpc.MethodDescriptor.<com.hashicorp.pong.PongData, com.hashicorp.pong.PongData>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.BIDI_STREAMING)
              .setFullMethodName(generateFullMethodName(
                  "pong.PongService", "ServerStream"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.hashicorp.pong.PongData.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.hashicorp.pong.PongData.getDefaultInstance()))
                  .setSchemaDescriptor(new PongServiceMethodDescriptorSupplier("ServerStream"))
                  .build();
          }
        }
     }
     return getServerStreamMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static PongServiceStub newStub(io.grpc.Channel channel) {
    return new PongServiceStub(channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static PongServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    return new PongServiceBlockingStub(channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static PongServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    return new PongServiceFutureStub(channel);
  }

  /**
   */
  public static abstract class PongServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public io.grpc.stub.StreamObserver<com.hashicorp.pong.PongData> clientStream(
        io.grpc.stub.StreamObserver<com.hashicorp.pong.PongData> responseObserver) {
      return asyncUnimplementedStreamingCall(getClientStreamMethod(), responseObserver);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<com.hashicorp.pong.PongData> serverStream(
        io.grpc.stub.StreamObserver<com.hashicorp.pong.PongData> responseObserver) {
      return asyncUnimplementedStreamingCall(getServerStreamMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getClientStreamMethod(),
            asyncBidiStreamingCall(
              new MethodHandlers<
                com.hashicorp.pong.PongData,
                com.hashicorp.pong.PongData>(
                  this, METHODID_CLIENT_STREAM)))
          .addMethod(
            getServerStreamMethod(),
            asyncBidiStreamingCall(
              new MethodHandlers<
                com.hashicorp.pong.PongData,
                com.hashicorp.pong.PongData>(
                  this, METHODID_SERVER_STREAM)))
          .build();
    }
  }

  /**
   */
  public static final class PongServiceStub extends io.grpc.stub.AbstractStub<PongServiceStub> {
    private PongServiceStub(io.grpc.Channel channel) {
      super(channel);
    }

    private PongServiceStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PongServiceStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new PongServiceStub(channel, callOptions);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<com.hashicorp.pong.PongData> clientStream(
        io.grpc.stub.StreamObserver<com.hashicorp.pong.PongData> responseObserver) {
      return asyncBidiStreamingCall(
          getChannel().newCall(getClientStreamMethod(), getCallOptions()), responseObserver);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<com.hashicorp.pong.PongData> serverStream(
        io.grpc.stub.StreamObserver<com.hashicorp.pong.PongData> responseObserver) {
      return asyncBidiStreamingCall(
          getChannel().newCall(getServerStreamMethod(), getCallOptions()), responseObserver);
    }
  }

  /**
   */
  public static final class PongServiceBlockingStub extends io.grpc.stub.AbstractStub<PongServiceBlockingStub> {
    private PongServiceBlockingStub(io.grpc.Channel channel) {
      super(channel);
    }

    private PongServiceBlockingStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PongServiceBlockingStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new PongServiceBlockingStub(channel, callOptions);
    }
  }

  /**
   */
  public static final class PongServiceFutureStub extends io.grpc.stub.AbstractStub<PongServiceFutureStub> {
    private PongServiceFutureStub(io.grpc.Channel channel) {
      super(channel);
    }

    private PongServiceFutureStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PongServiceFutureStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new PongServiceFutureStub(channel, callOptions);
    }
  }

  private static final int METHODID_CLIENT_STREAM = 0;
  private static final int METHODID_SERVER_STREAM = 1;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final PongServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(PongServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_CLIENT_STREAM:
          return (io.grpc.stub.StreamObserver<Req>) serviceImpl.clientStream(
              (io.grpc.stub.StreamObserver<com.hashicorp.pong.PongData>) responseObserver);
        case METHODID_SERVER_STREAM:
          return (io.grpc.stub.StreamObserver<Req>) serviceImpl.serverStream(
              (io.grpc.stub.StreamObserver<com.hashicorp.pong.PongData>) responseObserver);
        default:
          throw new AssertionError();
      }
    }
  }

  private static abstract class PongServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    PongServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.hashicorp.pong.Api.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("PongService");
    }
  }

  private static final class PongServiceFileDescriptorSupplier
      extends PongServiceBaseDescriptorSupplier {
    PongServiceFileDescriptorSupplier() {}
  }

  private static final class PongServiceMethodDescriptorSupplier
      extends PongServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    PongServiceMethodDescriptorSupplier(String methodName) {
      this.methodName = methodName;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.MethodDescriptor getMethodDescriptor() {
      return getServiceDescriptor().findMethodByName(methodName);
    }
  }

  private static volatile io.grpc.ServiceDescriptor serviceDescriptor;

  public static io.grpc.ServiceDescriptor getServiceDescriptor() {
    io.grpc.ServiceDescriptor result = serviceDescriptor;
    if (result == null) {
      synchronized (PongServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new PongServiceFileDescriptorSupplier())
              .addMethod(getClientStreamMethod())
              .addMethod(getServerStreamMethod())
              .build();
        }
      }
    }
    return result;
  }
}
