import os
import pathlib
import datetime


class TmpDir:
    """A temporary directory that is deleted when the object is destroyed."""

    tmpDirName = datetime.datetime.now().strftime("%Y-%m-%d")
    tmpFilePath = pathlib.Path(f"/home/ubuntu/NewsInbox/api/templates/{tmpDirName}/")

    def __init__(self):
        pathExists = os.path.exists(self.tmpFilePath)
        if not pathExists:
            os.makedirs(self.tmpFilePath)

    def path(self):
        return str(self.tmpFilePath) + "/"
