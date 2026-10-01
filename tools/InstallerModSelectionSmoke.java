import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import zombie.ZomboidFileSystem;
import zombie.core.znet.SteamUtils;
import zombie.gameStates.ChooseGameInfo;
import zombie.network.GameServer;

/** Real B42 filesystem and metadata parser; simulated installed-item list, no Steam or game loop. */
public final class InstallerModSelectionSmoke {
    private static Path mod(Path root, String id, String name) throws Exception {
        Files.createDirectories(root.resolve("common"));
        Files.createDirectories(root.resolve("42"));
        Files.writeString(root.resolve("42/mod.info"), "id=" + id + "\nname=" + name + "\n");
        return root.toAbsolutePath();
    }

    public static void main(String[] args) throws Exception {
        if (args.length != 1) throw new IllegalArgumentException("default or community");
        Path cache = Path.of("profile").toAbsolutePath();
        Path originalItem = Path.of("original-workshop").toAbsolutePath();
        Path community = mod(cache.resolve("mods/ZombieBuddy"), "ZombieBuddy", "ZombieBuddy Community");
        Path original = mod(originalItem.resolve("mods/ZombieBuddy"), "ZombieBuddy", "ZombieBuddy Original");
        mod(originalItem.resolve("mods/UnchangedConsumer"), "UnchangedConsumer", "Unchanged Consumer");
        Files.writeString(originalItem.resolve("mods/UnchangedConsumer/42/mod.info"),
                "id=UnchangedConsumer\nname=Unchanged Consumer\nrequire=ZombieBuddy\n");
        Files.createDirectories(Path.of("media"));
        ZomboidFileSystem.instance.setCacheDir(cache.toString());
        ZomboidFileSystem.instance.init();
        // Use the engine's server-side installed-folder input to avoid calling native Steam.
        // This exercises filesystem selection, not a dedicated-server or multiplayer session.
        GameServer.server = true;
        GameServer.workshopInstallFolders = new String[]{originalItem.toString()};
        var steamEnabled = SteamUtils.class.getDeclaredField("steamEnabled");
        steamEnabled.setAccessible(true);
        steamEnabled.setBoolean(null, true);
        if (args[0].equals("community")) ZomboidFileSystem.instance.setModFoldersOrder("mods,workshop,steam");
        ArrayList<String> folders = new ArrayList<>();
        ZomboidFileSystem.instance.getAllModFolders(folders);
        if (ChooseGameInfo.getModDetails("ZombieBuddy") == null) throw new AssertionError("Framework metadata missing");
        if (ChooseGameInfo.getModDetails("UnchangedConsumer") == null) throw new AssertionError("Consumer metadata missing");
        String selected = ZomboidFileSystem.instance.getModDir("ZombieBuddy");
        Path expected = args[0].equals("community") ? community : original;
        if (selected == null || !Path.of(selected).toAbsolutePath().startsWith(expected)) {
            throw new AssertionError("Expected " + expected + ", got " + selected + "; folders=" + folders);
        }
        if (ZomboidFileSystem.instance.getModDir("UnchangedConsumer") == null) {
            throw new AssertionError("Unchanged Workshop consumer missing");
        }
        System.out.println("PASS " + args[0] + ": selected=" + selected + "; unchanged Workshop consumer remains discoverable");
    }
}
