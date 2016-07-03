package ar.com.bestprice.buyitnow.dto;

import java.io.Serializable;

import ar.com.bestprice.buyitnow.Category;

/**
 * Created by ivan on 31/03/16.
 */
public class Item implements Serializable{

    private String id;
    private String description;
    private Float price;
    private String category = "";

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public Float getPrice() {
        return price;
    }

    public void setPrice(Float price) {
        this.price = price;
    }

    @Override
    public String toString() {
        return  description + "\t \t"+ price;
    }

    public Category getCategory() {

        if (category.isEmpty()){

            return Category.MERCADERIA;
        }else {

            return Category.valueOf(category.toUpperCase());
        }
    }

    public void setCategory(String category) {
        this.category = category;
    }
}
