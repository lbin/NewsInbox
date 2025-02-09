import os
import subprocess

from setuptools import find_packages, setup

version = '0.1.12'
package_name = 'newsinbox'
cwd = os.path.dirname(os.path.abspath(__file__))

sha = 'Unknown'
try:
    sha = subprocess.check_output(['git', 'rev-parse', 'HEAD'], cwd=cwd).decode('ascii').strip()
except Exception:
    pass


def write_version_file():
    version_path = os.path.join(cwd, 'newsinbox', 'version.py')
    with open(version_path, 'w') as f:
        f.write(f'__version__ = {version!r}\n')
        f.write(f'git_version = {repr(sha)}\n')


if __name__ == '__main__':
    print(f'Building wheel {package_name}-{version}')

    license = 'None'
    write_version_file()

    setup(
        name='newsinbox',
        version=version,
        url='',
        description='newsinbox',
        license=license,
        packages=find_packages(exclude=(
            'tests',
        )),
        zip_safe=False)