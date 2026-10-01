import java.util.Locale;
import me.zed_0xff.zombie_buddy.i18n.Messages;

/** Checks the final shaded JAR, not the source-tree resources. */
public class CatalogueSmoke {
    public static void main(String[] args) {
        Locale.setDefault(Locale.FRANCE);
        System.setProperty("zombiebuddy.language", "fr");
        if (!Messages.text("answer.yes").equals("Yes") || !Messages.text("column.updated").equals("Updated")) {
            throw new AssertionError("Approval text must remain English");
        }
        if (Messages.class.getResource("messages_fr.properties") != null) {
            throw new AssertionError("Unexpected French catalogue in the shaded JAR");
        }
        System.out.println("PASS: shaded English catalogue under non-English JVM locale; no French resource");
    }
}
