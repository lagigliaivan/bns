package ar.com.bestprice.buyitnow;

import android.content.Intent;
import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.view.View;
import android.widget.EditText;

import com.google.android.gms.common.api.CommonStatusCodes;


import ar.com.bestprice.buyitnow.dto.Item;

/**
 * Created by ivan on 01/04/16.
 */
public class AddItemActivity extends AppCompatActivity implements View.OnClickListener{

    @Override
    protected void onCreate(Bundle savedInstanceState) {

        super.onCreate(savedInstanceState);

        setContentView(R.layout.activity_add_item);

        String barcodeId = getIntent().getStringExtra("BarCode");

        EditText t = (EditText)findViewById(R.id.editText);

        t.setText(barcodeId);
    }

    @Override
    public void onClick(View v) {

        EditText id = (EditText)findViewById(R.id.editText);
        EditText description = (EditText)findViewById(R.id.editText2);
        EditText price = (EditText)findViewById(R.id.editText3);

        Item item = new Item();
        item.setId(id.getText().toString());
        item.setDescription(description.getText().toString());
        item.setPrice(Float.valueOf(price.getText().toString()));

        Intent data = new Intent();
        data.putExtra("Item",item);

        setResult(CommonStatusCodes.SUCCESS, data);
        finish();
    }
}
