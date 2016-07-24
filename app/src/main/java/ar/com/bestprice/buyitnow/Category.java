package ar.com.bestprice.buyitnow;

/**
 * Created by ivan on 03/07/16.
 */
public enum Category {

    MERCADERIA ("Mercaderia", 1, R.drawable.ic_mercaderia_24dp),
    DIVERSION   ("Diversion", 2, R.drawable.ic_diversion_24dp),
    IMPUESTOS   ("Impuestos", 3, R.drawable.ic_impuestos_24dp),
    SERVICIOS ("Servicios", 4, R.drawable.ic_servicios_24dp),
    SALUD ("Salud", 5, R.drawable.ic_salud_24dp),
    OTROS  ("Otros",6, R.drawable.ic_otros_24dp);


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
