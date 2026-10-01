package me.zed_0xff.zombie_buddy;

import java.io.InputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.security.cert.Certificate;
import java.security.cert.X509Certificate;
import java.util.jar.JarEntry;
import java.util.jar.JarFile;
import java.util.jar.Manifest;

/**
 * Retains the 2.x API and JAR-signature verification helpers.
 * Community previews use explicit manual updates: automatic replacement and
 * deferred .new writes are disabled at both public mutation entry points.
 */
public final class SelfUpdater {
    private SelfUpdater() {}

    private static String pendingNewVersion = null;

    public static String getNewVersion() {
        return pendingNewVersion;
    }

    public static String getExclusionReasonSuffix(Path jarPath) {
        if (jarPath == null) {
            return " (null not found)";
        }
        if (!Files.isRegularFile(jarPath)) {
            return " (" + jarPath.toAbsolutePath() + " not found)";
        }

        StringBuilder sb = new StringBuilder();
        Manifest manifest = getJarManifest(jarPath);
        if (manifest != null) {
            String manifestVersion = manifest.getMainAttributes().getValue("Implementation-Version");
            if (manifestVersion != null) {
                sb.append(" (version ").append(manifestVersion).append(")");
            }
        }

        Path currentJarPath = Utils.getCurrentJarPath();
        if (checkAndUpdateIfNewer(jarPath, currentJarPath, ZombieBuddy.getVersion(), Loader.g_verbosity)) {
            String newVer = getNewVersion();
            if (newVer != null) {
                sb.append(" -> updating to ").append(newVer);
            }
        }
        return sb.toString();
    }

    /** Preview updates are manual: never replace the fork with an upstream JAR. */
    public static void performUpdate(Path currentJarPath, Path newJarPath, String newVersion) {
        Logger.warn("ZombieBuddy Community preview uses manual updates; no JAR was replaced.");
    }

    /** Kept for binary compatibility with 2.x callers. */
    public static boolean checkAndUpdateIfNewer(Path jarPath, Path currentJarPath, String currentVersion, int verbosity) {
        return false;
    }

    public static Certificate[] verifyJarAndGetCerts(Path jarPath) throws Exception {
        Manifest mf = getJarManifest(jarPath);
        if (mf == null) {
            return null;
        }

        try (JarFile jar = new JarFile(jarPath.toFile(), true)) {
            Certificate[] signer = null;

            var entries = jar.entries();
            while (entries.hasMoreElements()) {
                JarEntry e = entries.nextElement();
                if (e.isDirectory()) continue;

                try (InputStream is = jar.getInputStream(e)) {
                    is.readAllBytes();
                }

                String name = e.getName();
                boolean isMeta = name.startsWith("META-INF/");
                Certificate[] certs = e.getCertificates();
                if (certs == null || certs.length == 0) {
                    if (!isMeta) {
                        throw new SecurityException("Unsigned entry: " + name);
                    }
                } else if (signer == null) {
                    signer = certs;
                }
            }

            if (signer == null) {
                throw new SecurityException("No signed entries found");
            }

            return signer;
        }
    }

    public static Manifest getJarManifest(Path jarPath) {
        if (jarPath == null) {
            return null;
        }
        try (JarFile jar = new JarFile(jarPath.toFile(), true)) {
            return jar.getManifest();
        } catch (Exception e) {
            Logger.error("Error getting JAR manifest: " + e);
            return null;
        }
    }

    private static byte[] getCertFingerprint(X509Certificate cert, int certNumber, int verbosity) {
        if (verbosity > 0) {
            Logger.info("  Certificate " + certNumber + ":");
            Logger.info("    Subject: " + cert.getSubjectX500Principal().getName());
            Logger.info("    Issuer: " + cert.getIssuerX500Principal().getName());
            Logger.info("    Serial Number: " + cert.getSerialNumber().toString(16).toUpperCase());
            Logger.info("    Valid From: " + cert.getNotBefore());
            Logger.info("    Valid Until: " + cert.getNotAfter());
        }

        try {
            byte[] sha256Bytes = Utils.sha256(cert.getEncoded());
            if (sha256Bytes != null && verbosity > 0) {
                Logger.info("    SHA-256 Fingerprint: " + Utils.bytesToHex(sha256Bytes, "%02X", ":"));
            }
            return sha256Bytes;
        } catch (Exception e) {
            Logger.error("    Error computing certificate fingerprints: " + e.getMessage());
        }
        return null;
    }
}
