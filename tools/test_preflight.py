import unittest
from preflight import inspect_options


class PreflightRegression(unittest.TestCase):
    def test_real_duplicate_installation_shape_is_blocked(self):
        agents, issues = inspect_options(["-Xmx3072m", "-agentlib:zbNative"],
            ["-agentlib:zbNative=patches_jar=AftermathLHCompat4221.jar:aftermathsystems.lhcompat"], True)
        self.assertEqual(len(agents), 2)
        self.assertEqual(len(issues), 2)

    def test_one_native_agent_with_unrelated_java_agent_is_accepted(self):
        _, issues = inspect_options(["-agentlib:zbNative"], ["-javaagent:AnotherMod.jar"], True)
        self.assertEqual(issues, [])

    def test_json_alone_is_not_claimed_complete(self):
        _, issues = inspect_options(["-javaagent:ZombieBuddy.jar"], [], False)
        self.assertTrue(issues)

    def test_one_quoted_java_agent_path_is_accepted(self):
        agents, issues = inspect_options([], ['-javaagent:"D:/Mod Files/ZombieBuddy.jar"=policy=deny-new'], True)
        self.assertEqual(len(agents), 1)
        self.assertEqual(issues, [])


if __name__ == "__main__":
    unittest.main()
