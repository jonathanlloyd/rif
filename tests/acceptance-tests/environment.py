# USE: behave -D BEHAVE_DEBUG_ON_ERROR         (to enable  debug-on-error)
import subprocess

BEHAVE_DEBUG_ON_ERROR = False


def before_all(context):
    setup_debug_on_error(context.config.userdata)

def setup_debug_on_error(userdata):
    global BEHAVE_DEBUG_ON_ERROR
    BEHAVE_DEBUG_ON_ERROR = userdata.getbool("BEHAVE_DEBUG_ON_ERROR")

def after_step(context, step):
    if BEHAVE_DEBUG_ON_ERROR and step.status == "failed":
        import ipdb
        ipdb.post_mortem(step.exc_traceback)


def before_feature(context, feature):
    if 'needs_echo_server' in feature.tags:
        print('starting echo server...')
        context.echo_server = subprocess.Popen(['/vol/tests/echo-server'], stdout=subprocess.PIPE, stderr=subprocess.STDOUT)

def after_feature(context, feature):
    if 'needs_echo_server' in feature.tags:
        print('stopping echo server')
        context.echo_server.terminate()
        context.echo_server.wait()
