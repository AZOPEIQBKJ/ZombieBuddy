import java.nio.file.Files;
import java.nio.file.Path;
import zombie.ZomboidFileSystem;
import zombie.core.znet.SteamWorkshopItem;

/** Read-only draft validation using the locally installed PZ 42.21 classes. No game loop or upload. */
public class ValidateWorkshopDraft {
    public static void main(String[] args) throws Exception {
        if (args.length != 2) throw new IllegalArgumentException("draftDirectory isolatedCacheDirectory");
        // Run from a disposable directory, never the game or user's active profile.
        Files.createDirectories(Path.of("media"));
        ZomboidFileSystem.instance.setCacheDir(args[1]);
        ZomboidFileSystem.instance.init();
        var item = new SteamWorkshopItem(args[0]);
        if (!item.readWorkshopTxt()) throw new AssertionError("Cannot read workshop.txt");
        if (item.getID() != null) throw new AssertionError("New draft must not target an existing item");
        if (item.getVisibilityInteger() != 2) throw new AssertionError("Draft must be private");
        if (!item.getDescription().contains("État actuel")) throw new AssertionError("French description lost");
        String error = item.validateContents();
        if (error != null) throw new AssertionError(error);
        System.out.println("PASS: PZ 42.21 metadata, content and PNG validation; private; no item ID; no upload/game loop");
    }
}
