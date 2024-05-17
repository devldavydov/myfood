// On window load.
$(function() {  
    $("#btnDelete").on("click", function() {
        if (!confirm("Удалить еду?"))
            return;
        
        apiPOST(FoodDeleteAPI, {key: foodKey})
        .done(function(resp) {
            enqueueNotification(NotifClassInfo, "Еда была удалена")
            window.location.replace(FoodListPage);
        })
        .fail(function(errMsg) {
            alert(errMsg);
        })
        .always(function() {
            console.log('always');
        });
    });

    $("#btnSet").on("click", function() {
        showNotification(NotifClassWarning, "Функционал в разработке!");
    })

    function loadFood() {
        apiGET(FoodGetAPI.replace(":key", foodKey))
        .done(function(resp) {
            if (resp.error !== "") {
                showEl($("#blockError"));    
                showNotification(NotifClassError, resp.error);
                return;
            }

            $("#name").val(resp.data.name);
            $("#brand").val(resp.data.brand);
            $("#cal100").val(resp.data.cal100.toFixed(2));
            $("#pfc").val(`${resp.data.prot100.toFixed(2)} / ${resp.data.fat100.toFixed(2)} / ${resp.data.carb100.toFixed(2)}`);
            $("#comment").val(resp.data.comment);

            showEl($("#blockFound"));
        })
        .fail(function(errMsg) {
            showEl($("#blockError"));
            showNotification(NotifClassError, errMsg);
        })
        .always(function() {
            hideEl($("#blockLoader"));
        });
    }

    loadFood();
});
