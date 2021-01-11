$(function () {

    const maxJiraFields = 3;
    const defaultJiraIssueField = 'fixVersions';
    const defaultJiraPattern = 'EDP_VERSION-EDP_COMPONENT';

    $($('section select.jiraIssueFields')[0]).selectize({
        create: true,
        sortField: {
            field: 'text',
            direction: 'asc'
        },
        dropdownParent: 'body'
    });

    !function preloadJiraMetadataFields() {
        $('select.jiraServer option').each(function (i, v) {
            let jiraServer = $(v).val();
            _sendGetRequest(true, `${$('input[id="basepath"]').val()}/api/v1/jira/${jiraServer}/metadata/fields`,
                function (fields) {
                    _setData(jiraServer, fields);
                }, function () {
                    console.error('an error has occurred while fetching Jira metadata fields')
                })
        });
    }();

    function _setData(jiraServer, fields) {
        let $advSettingsBlock = $('div.advanced-settings-block'),
            confData = $advSettingsBlock.attr('data-config');
        if (!!confData) {
            let tempConfData = $.parseJSON(confData);
            tempConfData[jiraServer] = fields;
            $advSettingsBlock.attr('data-config', JSON.stringify(tempConfData))
            return
        }
        $advSettingsBlock.attr('data-config', JSON.stringify({[jiraServer]: fields}))
    }

    $('#jiraServerToggle').change(function () {
        if ($(this).is(':checked')) {
            let jiraServer = $('select.jiraServer option:selected').val(),
                fields = $.parseJSON($('div.advanced-settings-block').attr('data-config'))[jiraServer];
            $.each(fields, function (i, item) {
                $('select.jiraIssueFields')[0].selectize.addOption({value: i, text: item});
            });
            _setDefaultValues();
            return
        }
    });

    function _setDefaultValues() {
        $('section select.jiraIssueFields')[0].selectize.setValue(defaultJiraIssueField);
        $('section div.jira-issue-metadata-row input.jiraPattern').val(defaultJiraPattern);
    }

    $('.add-jira-field').click(function () {
        if ($('.jiraIssueMetadata .jira-issue-metadata-row').length === maxJiraFields) {
            return
        }

        let $row = _getTemplate();
        $row.insertBefore($('.add-jira-field'));

        let removeButtonId = generateId();
        $row.find('.remove-jira-issue-metadata-row').addClass(removeButtonId)

        let selectId = generateId();
        $row.find('select').addClass(selectId)

        let $newSelect = $('section select.jiraIssueFields').last()
        $($newSelect).selectize({
            create: true,
            sortField: {
                field: 'text',
                direction: 'asc'
            },
            dropdownParent: 'body'
        });
        $row
            .removeClass('full-template')
            .removeClass('partial-template')
            .removeClass('hide-element');

        let $deleteBtn = $newSelect.parents('.jira-issue-metadata-row').find(`button.remove-jira-issue-metadata-row.${removeButtonId}`);
        $('.jiraIssueMetadata').on('click', `.remove-jira-issue-metadata-row.${removeButtonId}`, _removeJiraIssueRow.bind($deleteBtn));

        let $select = $newSelect.parents('.jira-issue-metadata-row').find(`select.${selectId}`);
        $('.jiraIssueMetadata').on('change', `.${selectId}`, toggleSelectOption.bind($select));

        let selectedValues = $('section select.jiraIssueFields').map(function () {
            return this.value
        }).get();
        let jiraServer = $('select.jiraServer option:selected').val(),
            fields = $.parseJSON($('div.advanced-settings-block').attr('data-config'))[jiraServer];
        $.each(fields, function (i, item) {
            if ($.inArray(i, selectedValues) !== -1) {
                return;
            }
            $newSelect[$newSelect.length - 1].selectize.addOption({value: i, text: item})
        });

        _tryToDisableButton();
    });

    function generateId() {
        return Date.now().toString(36) + Math.random().toString(36).substring(2);
    }

    function _getTemplate() {
        return $('.jiraIssueMetadata .jira-issue-metadata-row').length === 0
            ? $('.full-template.jira-issue-metadata-row').clone()
            : $('.partial-template.jira-issue-metadata-row').clone();
    }

    function _tryToDisableButton() {
        if ($('.jiraIssueMetadata .jira-issue-metadata-row').length === maxJiraFields) {
            $('.add-jira-field.circle.plus').addClass('disable').attr('disabled');
            return
        }
    }

    function _enableButton() {
        $('.add-jira-field.circle.plus').removeClass('disable').removeAttr('disabled');
    }

    $('.delete.remove-jira-issue-metadata-row').click(function () {
        _removeJiraIssueRow.bind(this)();
    });

    function _removeJiraIssueRow() {
        let $row = $(this).parents('.jira-issue-metadata-row'),
            $rows = $('section .jira-issue-metadata-row');
        if ($row.is(':first-child') && $rows.length >= 2) {
            let $fieldLabel = $('.jiraFieldNameLabel.template').clone(),
                $patternLabel = $('.jiraPatternLabel.template').clone();
            $fieldLabel.insertBefore($($rows[1]).find('select.jiraIssueFields'));
            $fieldLabel.removeClass('template').removeClass('hide-element');
            $patternLabel.insertBefore($($rows[1]).find('input.jiraPattern'));
            $patternLabel.removeClass('template').removeClass('hide-element');
        }

        toggleSelectOptions.bind($row.find('select.jiraIssueFields'))();
        $row.remove();
        _enableButton();
    }

    $('.jiraServer').change(function () {
        $('section div.jira-issue-metadata-row').remove();
        _enableButton();
    });

    function toggleSelectOption() {
        toggleSelectOptions.bind(this)();
    }

    $('.jiraIssueFields.jiraFieldName').change(function () {
        toggleSelectOptions.bind(this)();
    });

    function toggleSelectOptions() {
        let oldData;
        if (typeof $(this).attr('data-old') !== 'undefined') {
            oldData = $.parseJSON($(this).attr('data-old'));
        }

        let selectedOption = $(this).find('option:selected').val();
        $.each($('section .jira-issue-metadata-row select'), function () {
            if ($(this).find('option:selected').val() === selectedOption) {
                return
            }
            this.selectize.removeOption(selectedOption);

            if (typeof oldData !== 'undefined') {
                this.selectize.addOption({value: oldData.key, text: oldData.value});
            }
        });

        $(this).attr('data-old', JSON.stringify({
            key: $(this).val(),
            value: $(this).text()
        }));
    }
});