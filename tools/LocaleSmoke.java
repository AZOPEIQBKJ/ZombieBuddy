import java.nio.charset.StandardCharsets;
import java.nio.file.Path;
import me.zed_0xff.zombie_buddy.i18n.Messages;

/** Exercises the shaded JAR's UTF-8 catalogue and a child JVM with a different default locale. */
public class LocaleSmoke {
    public static void main(String[] args) throws Exception {
        if (args.length == 1 && args[0].equals("child")) {
            System.out.print(Messages.language() + "|" + Messages.text("answer.yes"));
            return;
        }
        if (args.length != 2) throw new IllegalArgumentException("agentJar smokeSource");
        if (!Messages.language().equals("fr") || !Messages.text("column.updated").equals("Mise à jour")) {
            throw new AssertionError("Shaded UTF-8 catalogue not loaded");
        }
        String executable = System.getProperty("os.name").toLowerCase().contains("win") ? "java.exe" : "java";
        Process child = new ProcessBuilder(Path.of(System.getProperty("java.home"), "bin", executable).toString(),
            "-Duser.language=en", Messages.childJvmArgument(), "-cp", args[0], args[1], "child")
            .redirectError(ProcessBuilder.Redirect.INHERIT).start();
        String response = new String(child.getInputStream().readAllBytes(), StandardCharsets.UTF_8);
        if (child.waitFor() != 0 || !response.equals("fr|Oui")) throw new AssertionError(response);
        System.out.println("PASS: shaded UTF-8 catalogue and normalized French locale propagated to English-default child JVM");
    }
}
