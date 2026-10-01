import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.LinkedList;
import java.util.List;
import me.zed_0xff.zombie_buddy.Agent;
import me.zed_0xff.zombie_buddy.Loader;

/** Compiled against unmodified 2.3.3, executed against the candidate. No game loop. */
public final class NativeListSmoke {
    public static void main(String[] args) throws Exception {
        Files.createDirectories(Path.of("media"));
        zombie.ZomboidFileSystem.instance.init();
        for (List<String> mods : List.of(new ArrayList<String>(), new LinkedList<String>(), List.<String>of())) {
            zombie.ZomboidFileSystem.instance.loadMods(mods);
        }
        // This descriptor is compiled against the baseline: a reflection-only test is insufficient.
        Loader.loadMods(new ArrayList<String>());
        if (args.length > 0 && args[0].equals("community")) {
            var reorder = Loader.class.getDeclaredMethod("autoFixModOrder", List.class);
            reorder.setAccessible(true);
            for (List<String> mods : List.of(new ArrayList<String>(), new LinkedList<String>())) {
                mods.addAll(List.of("first", "zdk", "second", "ZombieBuddy", "ZModUnbork", "third"));
                var expected = List.of("ZombieBuddy", "zdk", "ZModUnbork", "first", "second", "third");
                reorder.invoke(null, mods);
                if (!expected.equals(mods)) throw new AssertionError("in-place mod order changed");
                reorder.invoke(null, mods);
                if (!expected.equals(mods)) throw new AssertionError("mod order is not idempotent");
            }
        }
        if (!"deny-new".equals(Agent.arguments.get("policy"))) throw new AssertionError("policy changed");
        var policy = Loader.class.getDeclaredMethod("getPolicy");
        policy.setAccessible(true);
        if (!"deny-new".equals(policy.invoke(null))) throw new AssertionError("effective policy changed");
        if (Agent.arguments.containsKey("patches_jar")) throw new AssertionError("an adapter was supplied");
        System.out.println("PASS native List calls and legacy binary caller; approval policy unchanged");
    }
}
