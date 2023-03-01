-- replace agg_type_agg_id
BEGIN;
DROP INDEX eventstore.events@agg_type_agg_id;
COMMIT;

BEGIN;
CREATE INDEX agg_type_agg_id ON eventstore.events (
    instance_id
    , aggregate_type
    , aggregate_id
) STORING (
    event_type
    , aggregate_version
    , previous_aggregate_sequence
    , previous_aggregate_type_sequence
    , creation_date
    , event_data
    , editor_user
    , editor_service
    , resource_owner
);
COMMIT;

-- replace agg_type
BEGIN;
DROP INDEX eventstore.events@agg_type;
COMMIT;

BEGIN;
CREATE INDEX agg_type ON eventstore.events (
    instance_id
    , aggregate_type
    , event_sequence
) STORING (
    event_type
    , aggregate_id
    , aggregate_version
    , previous_aggregate_sequence
    , previous_aggregate_type_sequence
    , creation_date
    , event_data
    , editor_user
    , editor_service
    , resource_owner
);
COMMIT;

-- drop unused index
BEGIN;
DROP INDEX IF EXISTS eventstore.events@agg_type_seq;
COMMIT;

-- index to search event payload
BEGIN;
DROP INDEX IF EXISTS eventstore.events@event_search;
COMMIT;

BEGIN;
CREATE INVERTED INDEX event_search ON eventstore.events (
    instance_id
    , aggregate_type
    , event_type
    , event_data
);
COMMIT;