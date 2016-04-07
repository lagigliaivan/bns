package ar.com.bestprice.buyitnow;
import java.util.ArrayList;
import java.util.List;

import ar.com.bestprice.buyitnow.dto.Item;
import ar.com.bestprice.buyitnow.dto.Purchase;

public class Group {

    public String string;
    public final List<Item> children = new ArrayList<>();

    public Group(String string) {
        this.string = string;
    }

}
