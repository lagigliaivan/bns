package ar.com.bestprice.buyitnow.dto;

import java.util.ArrayList;
import java.util.Date;
import java.util.List;

/**
 * Created by ivan on 07/04/16.
 */
public class Purchase {

    String time;
    String shop;
    String id;
    List<Item> items = new ArrayList<>();

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

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

    public String getShop() {
        return shop;
    }

    public void setShop(String shop) {
        this.shop = shop;
    }

    public boolean isEmpty(){
        return items.isEmpty();
    }
}
