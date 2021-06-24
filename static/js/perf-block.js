$('#perfServerToggle').change(function () {
    let $perfBlockEl = $('div.perfServerBlock'),
        $dataSourceBlockEl = $('div.data-source-block');
    if ($(this).is(':checked')) {
        $perfBlockEl.removeClass('hide-element')
            .find('select[name="perfServer"]')
            .prop('disabled', false);
        $dataSourceBlockEl.removeClass('hide-element');
        displayButtonWhenPerfIsSelected();
        return
    }
    $perfBlockEl.addClass('hide-element')
        .find('select[name="perfServer"]')
        .prop('disabled', true);
    $dataSourceBlockEl.addClass('hide-element');
    displayButtonWhenPerfIsNotSelected();
});

function displayButtonWhenPerfIsSelected() {
    let $advBlockEl = $('div.card.advanced-settings-block'),
        $codebaseType = $advBlockEl.data('codebase-type');
    if ($codebaseType === 'application') {
        $('button.adv-setting-application-submit')
            .attr('disabled', false)
            .removeClass('hide-element')
            .attr('data-target', '#collapseDataSource');
        $('button.adv-setting-create-application')
            .attr('disabled', true)
            .addClass('hide-element');
    } else if ($codebaseType === 'autotests') {
        $('button.adv-setting-autotest-submit')
            .attr('disabled', false)
            .removeClass('hide-element')
            .attr('data-target', '#collapseDataSource');
        $('button.adv-setting-create-autotest')
            .attr('disabled', true)
            .addClass('hide-element');
    } else {
        $('button.adv-setting-library-submit')
            .attr('disabled', false)
            .removeClass('hide-element')
            .attr('data-target', '#collapseDataSource');
        $('button.adv-setting-create-library')
            .attr('disabled', true)
            .addClass('hide-element');
    }
}

function displayButtonWhenPerfIsNotSelected() {
    let $advBlockEl = $('div.card.advanced-settings-block'),
        vcsEnabled = $advBlockEl.data('vcs-enabled') === true,
        $codebaseType = $advBlockEl.data('codebase-type');
    if ($codebaseType === 'application') {
        if (vcsEnabled) {
            $('button.adv-setting-application-submit').attr('data-target', '#collapseVCS');
        } else {
            $('button.adv-setting-application-submit')
                .attr('disabled', true)
                .addClass('hide-element');
            $('button.adv-setting-create-application')
                .attr('disabled', false)
                .removeClass('hide-element');
        }
    } else if ($codebaseType === 'autotests') {
        if (vcsEnabled) {
            $('button.adv-setting-autotest-submit').attr('data-target', '#collapseVCS');
        } else {
            $('button.adv-setting-autotest-submit')
                .attr('disabled', true)
                .addClass('hide-element');
            $('button.adv-setting-create-autotest')
                .attr('disabled', false)
                .removeClass('hide-element');
        }
    } else {
        if (vcsEnabled) {
            $('button.adv-setting-library-submit').attr('data-target', '#collapseVCS');
        } else {
            $('button.adv-setting-library-submit')
                .attr('disabled', true)
                .addClass('hide-element');
            $('button.adv-setting-create-library')
                .attr('disabled', false)
                .removeClass('hide-element');
        }
    }
}

!function () {
    togglePerfDataSourcesView();
}();

$('#strategy').change(function () {
    togglePerfDataSourcesView();
});

$('select.ciTool').change(function () {
    toggleJenkinsDsView();
});

function togglePerfDataSourcesView() {
    let strategyType = $('#strategy').val(),
        $gitLabDsDivEl = $('.dataSources').find('div input[value="GitLab"]').parent('div'),
        importStrategy = 'import';

    if (importStrategy === strategyType.toLowerCase()) {
        $gitLabDsDivEl.removeClass('hide-element');
        toggleJenkinsDsView();
        return
    }
    $gitLabDsDivEl.addClass('hide-element');
}

function toggleJenkinsDsView() {
    let $jenkinsDsDivEl = $('.dataSources').find('div input[value="Jenkins"]').parent('div');
    if ($('select.ciTool option').filter(':selected').val() === 'GitLab CI') {
        $jenkinsDsDivEl.addClass('hide-element');
        return
    }
    $jenkinsDsDivEl.removeClass('hide-element');
}