package ar.com.bestprice.buyitnow;

import java.util.Enumeration;

/**
 * Created by ivan on 03/07/16.
 */
public enum Category {

    MERCADERIA ("Mercaderia", 1, R.drawable.shopping_cart_icon_48),
    DIVERSION   ("Diversion", 2, R.drawable.harlequin_red_icon_48),
    IMPUESTOS   ("Impuestos", 3, R.drawable.taxes_icon_48),
    SERVICIOS ("Servicios", 4, R.drawable.taxes_icon_48),
    SALUD ("Salud", 5, R.drawable.pill_icon_48),
    OTROS  ("Otros",6, R.drawable.joker_icon_48);


    private final String name;
    private final int id;
    private final int icon;

    Category(String name, int id, int drawable) {
        this.name = name;
        this.id = id;
        this.icon = drawable;
    }

    public String getName() { return name; }
    public int getId() { return id; }
    public int getIcon() {return icon;}
}
