import java.nio.file.Path;
import java.nio.file.Files;
import java.util.Objects;
import me.zed_0xff.zombie_buddy.ZombieBuddy;

/** Checks the actual packaged version and frozen consumer metadata, without a game loop. */
public final class ConsumerMetadataSmoke {
    public static void main(String[] args) throws Exception {
        Path versionDir = Path.of(args[0]).toAbsolutePath();
        if (!Files.isRegularFile(versionDir.resolve("mod.info")) ||
                !Files.isRegularFile(versionDir.resolve("media/java/AftermathSystemsLivingHordes.jar"))) {
            throw new AssertionError("Frozen consumer fixture is incomplete");
        }
        boolean expectAccepted = args[1].equals("accepted");
        Class<?> metadata = Class.forName("me.zed_0xff.zombie_buddy.JavaModInfo");
        var parse = metadata.getDeclaredMethod("parse", String.class);
        parse.setAccessible(true);
        Object result = parse.invoke(null, versionDir.toString());
        if ((result != null) != expectAccepted) {
            throw new AssertionError("Unexpected consumer result for " + ZombieBuddy.getVersion() + ": " + result);
        }
        if (result != null) {
            var jarPath = metadata.getDeclaredMethod("jarPath");
            jarPath.setAccessible(true);
            if (!Objects.equals(jarPath.invoke(result), versionDir.resolve("media/java/AftermathSystemsLivingHordes.jar"))) {
                throw new AssertionError("Wrong consumer JAR selected");
            }
        }
        System.out.println("PASS packaged " + ZombieBuddy.getVersion() + " consumer metadata " + args[1]);
    }
}
