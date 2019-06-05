function getUrlParameter(sParam) {
    let sPageURL = window.location.search.substring(1),
        sURLVariables = sPageURL.split('&'),
        sParameterName,
        i;

    for (i = 0; i < sURLVariables.length; i++) {
        sParameterName = sURLVariables[i].split('=');

        if (sParameterName[0] === sParam) {
            return sParameterName[1] === undefined ? true : decodeURIComponent(sParameterName[1]);
        }
    }
}

function _sendPostRequest(async, url, data, successCallback, failCallback) {
    $.ajax({
        url: url,
        contentType: "application/json",
        type: "POST",
        data: JSON.stringify(data),
        async: async,
        success: function (resp) {
            successCallback(resp);
        },
        error: function (resp) {
            failCallback(resp);
        }
    });
}

function isFieldValid(elementToValidate, regex) {
    let check = function (value) {
        return regex.test(value);
    };

    return !(!elementToValidate.val() || !check(elementToValidate.val()));
}

function blockIsNotValid($block) {
    $block.find('.card-header')
        .addClass('invalid')
        .removeClass('success')
        .addClass('error');
}

function blockIsValid($block) {
    $block.find('.card-header')
        .removeClass('invalid')
        .addClass('success')
        .removeClass('error');
}