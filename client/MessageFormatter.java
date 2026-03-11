package client;

import java.time.Instant;
import java.time.ZoneId;
import java.time.format.DateTimeFormatter;
import java.util.HashMap;
import java.util.Map;
import java.nio.file.Paths;

public class MessageFormatter {

    private static final DateTimeFormatter FORMATTER =
        DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss")
            .withZone(ZoneId.systemDefault());

    public static String formatTimestamp(long timestamp) {
        Instant instant = Instant.ofEpochMilli(timestamp);
        return FORMATTER.format(instant);
    }

    public static Map<String, Object> formatMessage(Map<String, Object> message) {
        Map<String, Object> formatted = new HashMap<>(message);

        if (message.containsKey("timestamp")) {
            long ts = (Long) message.get("timestamp");
            formatted.put("formatted_time", formatTimestamp(ts));
        }

        if (message.containsKey("content")) {
            String content = (String) message.get("content");
            formatted.put("sanitized_content", sanitize(content));
        }

        return formatted;
    }

    private static String sanitize(String input) {
        return input.replaceAll("<[^>]*>", "")
                   .replaceAll("&", "&amp;")
                   .replaceAll("\"", "&quot;");
    }

    public static String resolvePath(String pathStr) {
        return Paths.get(pathStr).toAbsolutePath().toString();
    }

    private static class MessageCache {
        private final Map<String, String> cache = new HashMap<>();

        public void put(String key, String value) {
            cache.put(key, value);
        }

        public String get(String key) {
            return cache.get(key);
        }
    }
}
