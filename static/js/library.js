$(function () {

    $('#languageSelection').on('change', function (e) {
        toggleJenkinsSlaveEl.call($('.advanced-settings-block'));
    });

    function toggleJenkinsSlaveEl() {
        _shouldJenkinsSlaveBeHidden() ?
            disableJenkinsSlaveEl.call(this) :
            enableJenkinsSlaveEl.call(this);
    }

    function disableJenkinsSlaveEl() {
        $(this).find('.jenkins-slave')
            .hide()
            .find('select.jenkinsSlave')
            .attr('disabled', true);
    }

    function enableJenkinsSlaveEl() {
        $(this).find('.jenkins-slave')
            .show()
            .find('select.jenkinsSlave')
            .attr('disabled', false);
    }

    function _shouldJenkinsSlaveBeHidden() {
        let $groovyPipeEl = $('.main-block div.card-body div.formSubsection-groovy-pipeline');
        return $groovyPipeEl.is(':visible') ?
            $groovyPipeEl.find('.groovy-pipeline-build-tools :selected').text() === 'none' : false;
    }

});
