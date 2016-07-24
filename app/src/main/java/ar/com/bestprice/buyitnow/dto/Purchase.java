package ar.com.bestprice.buyitnow.dto;

import java.util.ArrayList;
import java.util.Date;
import java.util.List;

/**
 * Created by ivan on 07/04/16.
 */
public class Purchase {

    String time;
    List<Item> items = new ArrayList<>();

    public List<Item> getItems() {

        if(items == null) {
            items = new ArrayList<>();
        }
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

    public void addItem(Item item) {

        if(item != null) {
            this.items.add(item);
        }
    }

    public void removeItem(Item item){
        items.remove(item);
    }
}
