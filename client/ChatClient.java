package client;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.WebSocket;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.CompletionStage;
import java.util.List;
import java.util.ArrayList;

public class ChatClient {
    private WebSocket webSocket;
    private final String serverUrl;
    private final ConcurrentLinkedQueue<String> messageQueue;
    private boolean connected;

    public ChatClient(String serverUrl) {
        this.serverUrl = serverUrl;
        this.messageQueue = new ConcurrentLinkedQueue<>();
        this.connected = false;
    }

    public void connect() {
        HttpClient client = HttpClient.newHttpClient();

        WebSocket.Builder builder = client.newWebSocketBuilder();

        this.webSocket = builder.buildAsync(
            URI.create(serverUrl),
            new WebSocket.Listener() {
                @Override
                public void onOpen(WebSocket webSocket) {
                    connected = true;
                    flushQueue();
                    WebSocket.Listener.super.onOpen(webSocket);
                }

                @Override
                public CompletionStage<?> onText(WebSocket webSocket, CharSequence data, boolean last) {
                    handleMessage(data.toString());
                    return WebSocket.Listener.super.onText(webSocket, data, last);
                }

                @Override
                public CompletionStage<?> onClose(WebSocket webSocket, int statusCode, String reason) {
                    connected = false;
                    return WebSocket.Listener.super.onClose(webSocket, statusCode, reason);
                }
            }
        ).join();
    }

    public void sendMessage(String content) {
        String message = String.format("{\"content\":\"%s\",\"timestamp\":%d}",
            content, System.currentTimeMillis());

        if (connected && webSocket != null) {
            webSocket.sendText(message, true);
        } else {
            messageQueue.offer(message);
        }
    }

    private void flushQueue() {
        while (!messageQueue.isEmpty()) {
            String message = messageQueue.poll();
            if (message != null) {
                webSocket.sendText(message, true);
            }
        }
    }

    private void handleMessage(String message) {
        System.out.println("Received: " + message);
    }

    public void disconnect() {
        if (webSocket != null) {
            webSocket.sendClose(WebSocket.NORMAL_CLOSURE, "Client closing");
        }
    }

    public List<String> getQueuedMessages() {
        return new ArrayList<>(messageQueue);
    }
}
