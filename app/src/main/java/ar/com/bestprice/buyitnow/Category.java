package ar.com.bestprice.buyitnow;

import java.util.Enumeration;

/**
 * Created by ivan on 03/07/16.
 */
public enum Category {

    MERCADERIA ("Mercaderia", 1, R.drawable.shop_cart_icon_32),
    DIVERSION   ("Diversion", 2, R.drawable.face_laugh_icon_32),
    IMPUESTOS   ("Impuestos", 3, R.drawable.taxes_icon_32),
    SERVICIOS ("Servicios", 4, R.drawable.payment_icon_32),
    SALUD ("Salud", 5, R.drawable.pill_icon_32),
    OTROS  ("Otros",6, R.drawable.basket_full_icon_32);


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
