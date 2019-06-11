create table if not exists cd_pipeline_service
(
 cd_pipeline_id integer not null
   constraint cd_pipeline_fk
     references cd_pipeline,
 service_id     integer not null
   constraint service_fk
     references service
);