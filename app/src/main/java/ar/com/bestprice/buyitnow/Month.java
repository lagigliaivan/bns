package ar.com.bestprice.buyitnow;

/**
 * Created by ivan on 03/07/16.
 */
public enum Month {


    JANUARY ("January", 1),
    FEBRUARY   ("February", 2),
    MARCH   ("March", 3),
    APRIL ("April", 4),
    MAY ("May", 5),
    JUNE  ("June", 6),
    JULY ("July", 7),
    AUGUST  ("August", 8),
    SEPTEMBER  ("September", 9),
    OCTOBER  ("October", 10),
    NOVEMBER  ("November", 11),
    DECEMBER  ("December", 12);


    private final String name;
    private final int position;


    Month(String name, int position) {
        this.name = name;
        this.position = position;
    }

    public String getName() { return name; }
    public int getPosition() { return position; }



}
