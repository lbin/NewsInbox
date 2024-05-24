import unittest
from unittest.mock import Mock, patch
from urllib.parse import urlparse

from newsinbox.servers.llm_common_post import sum4all


class TestSum4All(unittest.TestCase):
    @patch("requests.post")
    def test_sum4all(self, mock_post):
        # Mock the response from requests.post
        mock_response = Mock()
        mock_response.raise_for_status.return_value = None
        mock_response.json.return_value = {"choices": [{"message": {"content": "test content"}}]}
        mock_post.return_value = mock_response

        # Define a test config and content
        config = {
            "open_ai_api_base": "http://testbase.com",
            "open_ai_api_key": "testkey",
            "openai_chat_url": "http://testurl.com",
        }
        content = "test content"

        # Call the function with the test config and content
        result = sum4all(config, content)

        # Assert that requests.post was called with the correct arguments
        expected_headers = {"Authorization": "Bearer testkey", "Host": urlparse("http://testbase.com").netloc}
        mock_post.assert_called_once_with("http://testurl.com", headers=expected_headers, json=content, timeout=120)

        # Assert that the function returned the correct result
        self.assertEqual(result, "test content")


if __name__ == "__main__":
    unittest.main()
