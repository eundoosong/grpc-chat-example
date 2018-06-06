package io.grpc.chat;

import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.stub.StreamObserver;
import java.io.IOException;
import java.util.logging.Logger;
import java.util.logging.Level;
import java.util.UUID;

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

        public static String generateString() {
            return UUID.randomUUID().toString();
        }

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
        public StreamObserver<File> uploadFiles(final StreamObserver<FileIds> responseObserver) {
            return new StreamObserver<File>() {
                int fileCount;
                String id;
                @Override
                public void onNext(File file) {
                    fileCount++;
                    if(fileCount == 1) {
                        id = generateString();
                    }
                    logger.info("new file :" + file.getName() + ", " + file.getType());
                    logger.info("the number of files :" + fileCount);
                }

                @Override
                public void onError(Throwable t) {
                    logger.log(Level.WARNING, "Encountered error in recordRoute", t);
                }

                @Override
                public void onCompleted() {
                    logger.info("onComplete");
                    FileIds.Builder builder = FileIds.newBuilder();
                    for(int idx = 1; idx <= fileCount; idx++) {
                        builder.addId(id + "-" + idx);
                    }
                    responseObserver.onNext(builder.build());
                    responseObserver.onCompleted();
                }
            };
        }

        @Override
        public void downloadFiles(FileIds request, StreamObserver<File> responseObserver) {
            request.getIdCount();
            for(String id : request.getIdList()) {
                logger.info(id);
            }

            File.Builder builder = File.newBuilder();
            builder.setName("Test");
            builder.setType("plain/txt");
            builder.setLen(1);
            for(int i = 0; i < 10; i++) {
                responseObserver.onNext(builder.build());
            }
            responseObserver.onCompleted();
        }

        @Override
        public StreamObserver<File> convertFiles(StreamObserver<File> responseObserver) {
            return new StreamObserver<File>() {
                @Override
                public void onNext(File value) {
                    logger.info("onNext");
                    File.Builder builder = File.newBuilder();
                    builder.setName("Test");
                    builder.setType("plain/txt");
                    builder.setLen(1);
                    for(int i = 0; i < 2; i++) {
                        responseObserver.onNext(builder.build());
                    }
                }

                @Override
                public void onError(Throwable t) {
                    logger.log(Level.WARNING, "convertFiles cancelled");
                }

                @Override
                public void onCompleted() {
                    logger.info("onComplete");
                    responseObserver.onCompleted();
                }
            };
        }
    }
}
