import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.security.KeyFactory;
import java.security.MessageDigest;
import java.security.Signature;
import java.security.spec.X509EncodedKeySpec;
import java.util.HexFormat;

/** Standalone JDK verification; the trusted public key is supplied out of band. */
public final class VerifyZbs {
    public static void main(String[] args) throws Exception {
        if (args.length != 2 || !args[0].matches("[0-9a-fA-F]{64}")) {
            throw new IllegalArgumentException("Usage: VerifyZbs.java TRUSTED_PUBLIC_KEY_HEX JAR_PATH");
        }
        Path jar = Path.of(args[1]);
        var lines = Files.readAllLines(Path.of(args[1] + ".zbs"), StandardCharsets.UTF_8);
        if (lines.size() != 3 || !lines.get(0).equals("ZBS") ||
                !lines.get(1).matches("SteamID64:[0-9]{17}") ||
                !lines.get(2).matches("Signature:[0-9a-fA-F]{128}")) {
            throw new SecurityException("Invalid ZBS sidecar");
        }
        var hex = HexFormat.of();
        byte[] encoded = hex.parseHex("302a300506032b6570032100" + args[0]);
        var key = KeyFactory.getInstance("Ed25519").generatePublic(new X509EncodedKeySpec(encoded));
        String hash = hex.formatHex(MessageDigest.getInstance("SHA-256").digest(Files.readAllBytes(jar)));
        var verifier = Signature.getInstance("Ed25519");
        verifier.initVerify(key);
        verifier.update(("ZBS:" + lines.get(1).substring(10) + ":" + hash).getBytes(StandardCharsets.UTF_8));
        if (!verifier.verify(hex.parseHex(lines.get(2).substring(10)))) {
            throw new SecurityException("JAR signature does not match the trusted public key and content");
        }
        System.out.println("PASS ZBS signature " + lines.get(1) + " SHA256:" + hash);
    }
}
