package ar.com.bestprice.buyitnow;

import android.content.Intent;
import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.view.View;
import android.widget.EditText;

import com.google.android.gms.common.api.CommonStatusCodes;


import ar.com.bestprice.buyitnow.barcodereader.BarcodeCaptureActivity;
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

        EditText t = (EditText)findViewById(R.id.id_text);

        t.setText(barcodeId);
    }

    @Override
    public void onClick(View v) {

        EditText id = (EditText)findViewById(R.id.id_text);
        EditText description = (EditText)findViewById(R.id.description_text);
        EditText price = (EditText)findViewById(R.id.price_text);


        Item item = new Item();
        item.setId(id.getText().toString());
        item.setDescription(description.getText().toString());
        item.setPrice(Float.valueOf(price.getText().toString()));

        triggerActivity(item);

        finish();
    }

    private void triggerActivity(Item item){

        if (getIntent().getIntExtra(Constants.CALLING_ACTIVITY, 0) == Constants.NEW_PURCHASE){
            Intent data = new Intent();
            data.putExtra("Item", item);
            setResult(CommonStatusCodes.SUCCESS, data);

        } else {

            Intent intent = new Intent(this.getApplicationContext(), AddNewPurchaseActivity.class);
            intent.putExtra("Item", item);
            startActivity(intent);
        }
    }
}
