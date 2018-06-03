
package io.grpc.chat;

import static java.util.concurrent.TimeUnit.NANOSECONDS;

import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.stub.StreamObserver;
import java.io.IOException;
import java.util.logging.Logger;
import java.util.logging.Level;

public class ChatServer {
    private static final Logger logger = Logger.getLogger(ChatServer.class.getName());

    private Server server;

    private void start() throws IOException {
        int port = 50051;
        server = ServerBuilder.forPort(port)
                .addService(new ChatServiceImpl())
                .build()
                .start();
        logger.info("Server started, listening on " + port);
        Runtime.getRuntime().addShutdownHook(new Thread() {
            @Override
            public void run() {
                System.err.println("*** shutting down gRPC server since JVM is shutting down");
                ChatServer.this.stop();
                System.err.println("*** server shut down");
            }
        });
    }

    private void stop() {
        if(server != null) {
            server.shutdown();
        }
    }

    private void blockUntilShutdown() throws InterruptedException {
        if(server != null) {
            server.awaitTermination();
        }
    }

    public static void main(String[] args) throws IOException, InterruptedException {
        final ChatServer server = new ChatServer();
        server.start();
        server.blockUntilShutdown();
    }

    static class ChatServiceImpl extends ChatServiceGrpc.ChatServiceImplBase {

        @Override
        public void sendMessage(Message req, StreamObserver<Message> responseObserver) {
            Message reply = Message.newBuilder()
                    .setId("Hello " + req.getId())
                    .setText("Nice to meet you, " + req.getText())
                    .build();
            responseObserver.onNext(reply);
            responseObserver.onCompleted();
        }

        @Override
        public StreamObserver<File> sendFile(final StreamObserver<Url> responseObserver) {
            return new StreamObserver<File>() {
                @Override
                public void onNext(File file) {
                    file.getData();
                    logger.info("new file :" + file.getName() + ", " + file.getType());
                }

                @Override
                public void onError(Throwable t) {
                    logger.log(Level.WARNING, "Encountered error in recordRoute", t);
                }

                @Override
                public void onCompleted() {
                    logger.info("onComplete");
                    responseObserver.onNext(Url.newBuilder().setUrl("localhost:8080").build());
                    responseObserver.onCompleted();
                }
            };
        }

    }
}
