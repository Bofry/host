package app

import "testing"

func TestWithSessionStateManager(t *testing.T) {
	opt := WithSessionStateManager(nil)

	expected := APP_BUILDING_OPT_SESSION_STATE_MANAGER
	if expected != opt.typeName() {
		t.Errorf("assert typeName():: expected: %+v, got, %+v", expected, opt.typeName())
	}
}

func TestWithProtocolResolver(t *testing.T) {
	opt := WithProtocolResolver(nil)

	expected := APP_BUILDING_OPT_PROTOCOL_RESOLVER
	if expected != opt.typeName() {
		t.Errorf("assert typeName():: expected: %+v, got, %+v", expected, opt.typeName())
	}
}

func TestWithProtocolEmitter(t *testing.T) {
	opt := WithProtocolEmitter(nil)

	expected := APP_BUILDING_OPT_PROTOCOL_EMITTER
	if expected != opt.typeName() {
		t.Errorf("assert typeName():: expected: %+v, got, %+v", expected, opt.typeName())
	}
}

func TestWithInvalidMessageHandler(t *testing.T) {
	opt := WithInvalidMessageHandler(nil)

	expected := APP_BUILDING_OPT_INVALID_MESSAGE_HANDLER
	if expected != opt.typeName() {
		t.Errorf("assert typeName():: expected: %+v, got, %+v", expected, opt.typeName())
	}
}

func TestWithInvalidEventHandler(t *testing.T) {
	opt := WithInvalidEventHandler(nil)

	expected := APP_BUILDING_OPT_INVALID_EVENT_HANDLER
	if expected != opt.typeName() {
		t.Errorf("assert typeName():: expected: %+v, got, %+v", expected, opt.typeName())
	}
}

func TestWithLoggerOutput(t *testing.T) {
	opt := WithLoggerOutput(nil)

	expected := APP_BUILDING_OPT_LOGGER_OUTPUT
	if expected != opt.typeName() {
		t.Errorf("assert typeName():: expected: %+v, got, %+v", expected, opt.typeName())
	}
}

func TestWithLoggerFlags(t *testing.T) {
	opt := WithLoggerFlags(0)

	expected := APP_BUILDING_OPT_LOGGER_FLAGS
	if expected != opt.typeName() {
		t.Errorf("assert typeName():: expected: %+v, got, %+v", expected, opt.typeName())
	}
}

func TestWithTracerProvider(t *testing.T) {
	opt := WithTracerProvider(nil)

	expected := APP_BUILDING_OPT_TRACER_PROVIDER
	if expected != opt.typeName() {
		t.Errorf("assert typeName():: expected: %+v, got, %+v", expected, opt.typeName())
	}
}

func TestWithTextMapPropagator(t *testing.T) {
	opt := WithTextMapPropagator(nil)

	expected := APP_BUILDING_OPT_TEXT_MAP_PROPAGATOR
	if expected != opt.typeName() {
		t.Errorf("assert typeName():: expected: %+v, got, %+v", expected, opt.typeName())
	}
}

func TestWithEventClient(t *testing.T) {
	opt := WithEventClient(nil)

	expected := APP_BUILDING_OPT_EVENT_CLIENT
	if expected != opt.typeName() {
		t.Errorf("assert typeName():: expected: %+v, got, %+v", expected, opt.typeName())
	}
}

func TestWithMessageRouter(t *testing.T) {
	opt := WithMessageRouter(nil)

	expected := APP_BUILDING_OPT_MESSAGE_ROUTER
	if expected != opt.typeName() {
		t.Errorf("assert typeName():: expected: %+v, got, %+v", expected, opt.typeName())
	}
}

func TestWithEventRouter(t *testing.T) {
	opt := WithEventRouter(nil)

	expected := APP_BUILDING_OPT_EVENT_ROUTER
	if expected != opt.typeName() {
		t.Errorf("assert typeName():: expected: %+v, got, %+v", expected, opt.typeName())
	}
}

func TestWithDefaultMessageHandler(t *testing.T) {
	opt := WithDefaultMessageHandler(nil)

	expected := APP_BUILDING_OPT_DEFAULT_MESSAGE_HANDLER
	if expected != opt.typeName() {
		t.Errorf("assert typeName():: expected: %+v, got, %+v", expected, opt.typeName())
	}
}

func TestWithDefaultEventHandler(t *testing.T) {
	opt := WithDefaultEventHandler(nil)

	expected := APP_BUILDING_OPT_DEFAULT_EVENT_HANDLER
	if expected != opt.typeName() {
		t.Errorf("assert typeName():: expected: %+v, got, %+v", expected, opt.typeName())
	}
}

func TestWithErrorHandler(t *testing.T) {
	opt := WithErrorHandler(nil)

	expected := APP_BUILDING_OPT_ERROR_HANDLER
	if expected != opt.typeName() {
		t.Errorf("assert typeName():: expected: %+v, got, %+v", expected, opt.typeName())
	}
}
