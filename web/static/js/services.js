var gorunnerServices = angular.module('gorunnerServices', ['ngResource']);

gorunnerServices.factory('Job', ['$resource', function($resource){
  return $resource('/jobs/:id', {}, {
    update: { method: "PUT", params: {id: '@id'}},
    addTask: { method: "POST", url: '/jobs/:id/tasks', params: {id: '@id'}},
    removeTask: { method: "DELETE", url: '/jobs/:id/tasks/:tidx',  params: {id: '@id', tid: '@tidx'}},
    addTrigger: { method: "POST", url: '/jobs/:id/triggers', params: {id: '@id'}},
    removeTrigger: { method: "DELETE", url: '/jobs/:id/triggers/:trigger', params: {id: '@id', trigger: '@trigger'}}
  })
}]);

gorunnerServices.factory('Task', ['$resource', function($resource){
  return $resource('/tasks/:id', {}, {
    update: { method: "PUT" , params: {id: '@id'}},
    jobs: { method: "GET", url: 'tasks/:id/jobs', params: {id: '@id'}, isArray: true}
  })
}]);

gorunnerServices.factory('Trigger', ['$resource', function($resource){
  return $resource('/triggers/:id', {}, {
    update: { method: "PUT", params: {id: '@id'}},
    listJobs: { method: "GET", url: '/triggers/:id/jobs', params: {id: '@id'}, isArray: true}
  })
}]);

gorunnerServices.factory('Run', ['$resource', function($resource) {
  return $resource('/runs/:id', {}, {})
}]);
