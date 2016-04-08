package ar.com.bestprice.buyitnow.dto;

import java.util.Date;
import java.util.List;

/**
 * Created by ivan on 07/04/16.
 */
public class Purchase {

    String time;
    List<Item> items;

    public List<Item> getItems() {
        return items;
    }

    public void setItems(List<Item> items) {
        this.items = items;
    }

    public String getTime() {
        return time;
    }

    public void setTime(String time) {
        this.time = time;
    }
}
