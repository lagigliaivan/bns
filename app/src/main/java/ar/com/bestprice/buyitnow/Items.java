package ar.com.bestprice.buyitnow;

import java.util.ArrayList;
import java.util.List;

/**
 * Created by ivan on 01/04/16.
 */
public class Items {

    private List<Item> items = new ArrayList<>();


    public List getItems() {
        return items;
    }

    public void setItems(List<Item> items) {
        this.items = items;
    }
}
