"""Step definitions for acceptance tests"""

@given(u'a .rif file is on disk that describes a GET request from that server')
def step_impl(context):
    raise NotImplementedError(u'STEP: Given a .rif file is on disk that describes a GET request from that server')

@when(u'the user runs RIF on that file')
def step_impl(context):
    raise NotImplementedError(u'STEP: When the user runs RIF on that file')

@then(u'RIF should return an echo of the request it made')
def step_impl(context):
    raise NotImplementedError(u'STEP: Then RIF should return an echo of the request it made')
