$(function () {

    !function () {
        var $menuBtn = $('.js-toggle-menu-button');
        var $menu = $('.js-aside-menu');

        $menuBtn.on('click', function () {
            $(this).toggleClass('collapsed');
            $menu.toggleClass('active');
        });
    }();

    !function () {
        var $subFormWrap = $('#subFormWrapper');
        var $subForms = $('.js-form-subsection');

        $subFormWrap.on('change', function (e) {
            var target = $(e.target).data('target');
            $subForms.hide();
            $(target).show();
        });
    }();

    !function () {
        $.each($('.route, .dataBase .card-body').find('input, select'), function () {
            if ($(this).is('input')) {
                $(this).attr('readonly', true);
            } else if ($(this).is('select')) {
                $(this).attr('disabled', true);
            }
        });
    }();

    $(document).ready(function () {
        _sendGetRequest(true,'/api/v1/storage-class',
            function (storageClasses) {
                var $select = $('#dbPersistentStorage');

                $.each(storageClasses, function () {
                    $select.append('<option value="' + this.toString() + '">' + this.toString() + '</option>');
                });
            }, function (resp) {
                console.log(resp);
            })
    });

    $('.tooltip-icon').tooltip();

    $('#strategy').change(function () {
        if (this.value === 'clone') {
            $(".repo-url, .private-repo").removeClass('hide-element');
            if ($('#isRepoPrivate').is(':checked')) {
                $('.repoLogin, .repoPassword').removeClass('hide-element');
            }
        } else {
            $(".repo-url, .private-repo").addClass('hide-element');
            $('.repoLogin, .repoPassword').addClass('hide-element');
        }
    });

    $('#isRepoPrivate').change(function () {
        if ($(this).is(':checked')) {
            $('.repoLogin, .repoPassword').removeClass('hide-element');
        } else {
            $('.repoLogin, .repoPassword').addClass('hide-element');
            $('.repoLogin, .repoPassword').find('.invalid-feedback').hide();
            $('.repoLogin, .repoPassword').find('input').removeClass('is-invalid');

        }
    });

    /*hides and removes validate classes from optional fields (all route and database)*/
    $('#needRoute, #needDb').change(function () {
        toggleFields.bind(this)();

        $.each($(this).closest('.card-body').find('div.form-group:not(.hide-element) input'), function () {
            $(this).removeClass('is-invalid');
            $(this).next('.invalid-feedback').hide();
        });
    });

    $('.java-build-tools').change(function () {
        if (this.value === 'Maven') {
            $('#multiModule').attr("disabled", false);
            $('.multi-module').removeClass('hide-element');
        } else {
            $('.multi-module').addClass('hide-element');
            $('#multiModule').attr("disabled", true);
        }
    });

    /*disables other build tools that doesnt related to selected language*/
    /*todo think about better approach*/
    $('#collapseTwo .card-body .form__input-wrapper #subFormWrapper .form__radio-btn-wrapper .form__radio-btn').click(function () {
        $(this).parents('.card-body').find('.appLangError').hide();
        $('.framework').prop('checked', false);
        $('.multi-module').addClass('hide-element');
        $('.java-build-tools').val("Gradle");
        $('#multiModule').prop( "checked", false ).attr("disabled", true);
        var $target = $(this).data('target');
        $target = $target.substring(1, $target.length);

        $.each($(this).parents('.card-body').find('.form-group .form-subsection'), function () {
            if (!$(this).hasClass($target)) {
                $(this).find('select').attr('disabled', true);
            } else {
                $(this).find('select').removeAttr('disabled');
            }
        });
    });

    $('#btn-modal-continue').click(function () {
        $('#createAppForm').submit();
        $('#confirmationPopup').modal('hide');
        $(".window-table-body").remove();
    });

    $("#btn-cross-close, #btn-modal-close").click(function () {
        $(".window-table-body").remove();
    });

});

function toggleFields() {
    var toggleInputs = function (bool) {
        $.each($(this).closest('.card-body').find('input, select'), function () {
            if ($(this).is('input')) {
                $(this).attr('readonly', bool);
            } else if ($(this).is('select')) {
                $(this).attr('disabled', bool);
            }
        })
    }.bind(this);

    if ($(this).is(":checked")) {
        toggleInputs(false);
    } else {
        toggleInputs(true);
    }
}

function _sendGetRequest(async, url, successCallback, failCallback) {
    $.ajax({
        url: url,
        contentType: "application/json",
        async: async,
        success: function (resp) {
            successCallback(resp);
        },
        error: function (resp) {
            failCallback(resp);
        },
    });
}
