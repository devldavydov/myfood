$(function() {
    function loadFoodList() {
        apiGET(FoodListAPI)
        .done(function(resp) {
            if (resp.error !== "") {
                showNotification(NotifClassError, resp.error);
                return;
            }

            let tbl = $("#tblFood")
            for (let food of resp.data) {
                tbl.append(`
                <tr>
                    <td class="align-middle col-4">${food.name}</td>
                    <td class="align-middle col-2">${food.brand}</td>
                    <td class="align-middle col-1">${food.cal100}</td>
                    <td class="align-middle col-4">${food.comment}</td>
                    <td class="align-middle text-center col-1">
                        <a
                            href="/food/${food.key}"
                            class="btn btn-primary"><i class="bi bi-box-arrow-right"></i></a>
                    </td>
                </tr>     
                `)
            }       
        })
        .fail(function(errMsg) {
            showNotification(NotifClassError, errMsg);
        })
        .always(function() {
            hideEl($("#blockLoader"));
            showEl($("#blockList"));
        });
    }

    loadFoodList();
})