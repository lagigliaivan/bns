package ar.com.bestprice.buyitnow;

/**
 * Created by ivan on 03/07/16.
 */
public enum Month {


    //DO NOT CHANGE THE ORDER. This will affect the way the purchases are shown on the screen
    JANUARY ("January"),
    FEBRUARY ("February"),
    MARCH ("March"),
    APRIL ("April"),
    MAY ("May"),
    JUNE ("June"),
    JULY ("July"),
    AUGUST ("August"),
    SEPTEMBER ("September"),
    OCTOBER  ("October"),
    NOVEMBER ("November"),
    DECEMBER ("December");


    private final String name;

    Month(String name) {
        this.name = name;
    }

    public String getName() { return name; }
}
