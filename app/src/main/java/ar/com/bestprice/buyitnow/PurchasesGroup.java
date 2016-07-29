package ar.com.bestprice.buyitnow;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import ar.com.bestprice.buyitnow.dto.Item;
import ar.com.bestprice.buyitnow.dto.Purchase;

public class PurchasesGroup {

    private Month month;
    private float purchasesTotalPrice = 0;
    private Map<String, Purchase> time_purchases = new HashMap<>();
    private final List<Item> children = new ArrayList<>();

    public void addItem(Item item){
        children.add(item);
        purchasesTotalPrice += item.getPrice();
    }

    public List<Item> getChildren() {
        return children;
    }

    public String getString() {
        return month.toString();
    }

    public Month getMonth() {
        return month;
    }

    public PurchasesGroup(Month month) {
        this.month = month;
    }


    public float getPurchasesTotalPrice() {
        return purchasesTotalPrice;
    }


    public void addPurchase(Purchase purchase){
        time_purchases.put(purchase.getTime(), purchase);
    }

    public Purchase getPurchase(String time){
        return time_purchases.get(time);
    }

    public void removeItemAt(int childPosition) {

        Item item = children.remove(childPosition);
        Purchase purchase = time_purchases.get(item.getTime());
        purchase.removeItem(item);

        purchasesTotalPrice -= item.getPrice();
    }

    public Item getItemAt(int childPosition) {
        return children.get(childPosition);
    }
}
